package redshift

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/redshiftdataapiservice"
	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/sunker/redshift-datasource/pkg/models"
)

// SampleDatasource is an example datasource used to scaffold
// new datasource plugins with an backend.
type RedshiftDatasource struct {
	// The instance manager can help with lifecycle management
	// of datasource instances in plugins. It's not a requirements
	// but a best practice that we recommend that you follow.
	im instancemgmt.InstanceManager
}

// newDatasource returns datasource.ServeOpts.
func NewDatasource() datasource.ServeOpts {
	im := datasource.NewInstanceManager(newDataSourceInstance)

	ds := &RedshiftDatasource{
		im: im,
	}

	return datasource.ServeOpts{
		QueryDataHandler:   ds,
		CheckHealthHandler: ds,
	}
}


// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifer).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (td *RedshiftDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	log.DefaultLogger.Info("QueryData", "request", req)

	// create response struct
	response := backend.NewQueryDataResponse()

	// loop over queries and execute them individually.
	for _, q := range req.Queries {
		res := td.query(ctx, q, req.PluginContext)

		// save the response in a hashmap
		// based on with RefID as identifier
		response.Responses[q.RefID] = res
	}

	return response, nil
}

type queryModel struct {
	RawSQL  string `json:"rawSql,omitempty"`
	// Not from JSON
	Interval      time.Duration     `json:"-"`
	TimeRange     backend.TimeRange `json:"-"`
	MaxDataPoints int64             `json:"-"`
}

func (ds *RedshiftDatasource) query(ctx context.Context, query backend.DataQuery, pluginContext backend.PluginContext) backend.DataResponse {
	// Unmarshal the json into our queryModel
	var qm queryModel
	response := backend.DataResponse{}

	response.Error = json.Unmarshal(query.JSON, &qm)
	if response.Error != nil {
		return response
	}

	s, err := ds.getInstance(pluginContext)
	if err != nil {
		response.Error = err
		return response
	}

	client, err := s.ClientFactory("us-east-2")
	if err != nil {
		response.Error = err
		return response
	}

	input := &redshiftdataapiservice.ExecuteStatementInput{
		DbUser: aws.String("cloud-datasources"),
		Sql	: aws.String(qm.RawSQL),
		Database: aws.String("dev"),
		ClusterIdentifier: aws.String("redshift-cluster-grafana"),
	}
	res, err := client.ExecuteStatement(input)

	statement, err := client.DescribeStatement(&redshiftdataapiservice.DescribeStatementInput{
		Id: res.Id,
	})
	log.DefaultLogger.Info("QueryData", "statementres", statement.GoString())

	sql, err := client.GetStatementResult(&redshiftdataapiservice.GetStatementResultInput{
		Id: res.Id,
	})

	

	log.DefaultLogger.Info("QueryData", "sqlres", sql.GoString())

	// create data frame response
	frame := data.NewFrame("response")

	// add the time dimension
	frame.Fields = append(frame.Fields,
		data.NewField("time", nil, []time.Time{query.TimeRange.From, query.TimeRange.To}),
	)

	// add values
	frame.Fields = append(frame.Fields,
		data.NewField("values", nil, []int64{10, 20}),
	)

	// add the frames to the response
	response.Frames = append(response.Frames, frame)

	return response
}


func (ds *RedshiftDatasource) getInstance(ctx backend.PluginContext) (*instanceSettings, error) {
	s, err := ds.im.Get(ctx)
	if err != nil {
		return nil, err
	}
	return s.(*instanceSettings), nil
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (ds *RedshiftDatasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	s, err := ds.getInstance(req.PluginContext)
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: err.Error(),
		}, nil
	}

	client, err := s.ClientFactory("us-east-2")

	input := &redshiftdataapiservice.ExecuteStatementInput{
		DbUser: aws.String("cloud-datasources"),
		Sql	: aws.String("select * from public.category"),
		Database: aws.String("dev"),
		ClusterIdentifier: aws.String("redshift-cluster-grafana"),
	}
	res, err := client.ExecuteStatement(input)

	statement, err := client.DescribeStatement(&redshiftdataapiservice.DescribeStatementInput{
		Id: res.Id,
	})
	log.DefaultLogger.Info("QueryData", "statementres", statement.GoString())

	sql, err := client.GetStatementResult(&redshiftdataapiservice.GetStatementResultInput{
		Id: res.Id,
	})

	log.DefaultLogger.Info("QueryData", "sqlres", sql.GoString())

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "ok",
	}, nil
}

type clientFactoryFunc func(region string) (client *redshiftdataapiservice.RedshiftDataAPIService, err error)

type instanceSettings struct {
	ClientFactory   clientFactoryFunc
	Settings models.AWSRedshiftDataSourceSetting
}

func newDataSourceInstance(s backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	settings := models.AWSRedshiftDataSourceSetting{}
	err := settings.Load(s)
	if err != nil {
		return nil, fmt.Errorf("error reading settings: %s", err.Error())
	}
	sessions := awsds.NewSessionCache()

	return &instanceSettings{
		Settings: settings,
		ClientFactory: func(region string) (client *redshiftdataapiservice.RedshiftDataAPIService, err error) {
			session, err := sessions.GetSession(region, settings.AWSDatasourceSettings)
			if err != nil {
				return nil, err
			}
			return redshiftdataapiservice.New(session), nil
		},
	}, nil
}

func (s *instanceSettings) Dispose() {
	// Called before creatinga a new instance to allow plugin authors
	// to cleanup.
}
