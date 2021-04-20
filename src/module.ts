import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './datasource';
import { ConfigEditor } from './ConfigEditor';
import { QueryEditor } from './QueryEditor';
import { RedShiftDataSourceJsonData, RedShiftDataSourceSecureJsonData, RedshiftQuery } from './types';

export const plugin = new DataSourcePlugin<
  DataSource,
  RedshiftQuery,
  RedShiftDataSourceJsonData,
  RedShiftDataSourceSecureJsonData
>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
