import { DataSourceInstanceSettings } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { RedShiftDataSourceJsonData, RedShiftQuery } from './types';

export class DataSource extends DataSourceWithBackend<RedShiftQuery, RedShiftDataSourceJsonData> {
  constructor(instanceSettings: DataSourceInstanceSettings<RedShiftDataSourceJsonData>) {
    super(instanceSettings);
  }
}
