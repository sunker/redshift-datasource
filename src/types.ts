import { DataQuery } from '@grafana/data';
import { AwsAuthDataSourceJsonData, AwsAuthDataSourceSecureJsonData } from '@grafana/aws-sdk';

export interface RedShiftQuery extends DataQuery {
  queryText?: string;
  constant: number;
}

export const defaultQuery: Partial<RedShiftQuery> = {
  constant: 6.5,
};

/**
 * These are options configured for each DataSource instance
 */

export interface RedShiftDataSourceJsonData extends AwsAuthDataSourceJsonData {}

export interface RedShiftDataSourceSecureJsonData extends AwsAuthDataSourceSecureJsonData {}
/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface MySecureJsonData {
  apiKey?: string;
}
