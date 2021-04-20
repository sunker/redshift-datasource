import { DataQuery } from '@grafana/data';
import { AwsAuthDataSourceJsonData, AwsAuthDataSourceSecureJsonData } from '@grafana/aws-sdk';

export enum FormatOptions {
  TimeSeries,
  Table,
}

export interface RedshiftQuery extends DataQuery {
  rawSql: string;
  format: FormatOptions;
}

export const defaultQuery: Partial<RedshiftQuery> = {};

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
