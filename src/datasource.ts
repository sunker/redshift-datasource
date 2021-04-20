import { DataSourceInstanceSettings } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { RedShiftDataSourceJsonData, RedshiftQuery } from './types';

export class DataSource extends DataSourceWithBackend<RedshiftQuery, RedShiftDataSourceJsonData> {
  constructor(instanceSettings: DataSourceInstanceSettings<RedShiftDataSourceJsonData>) {
    super(instanceSettings);
  }
}
