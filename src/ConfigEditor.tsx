import React, { PureComponent } from 'react';
import { ConnectionConfig } from '@grafana/aws-sdk';
import { RedShiftDataSourceJsonData, RedShiftDataSourceSecureJsonData } from 'types';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';

export type Props = DataSourcePluginOptionsEditorProps<RedShiftDataSourceJsonData, RedShiftDataSourceSecureJsonData>;

export class ConfigEditor extends PureComponent<Props> {
  constructor(props: Props) {
    super(props);
    this.state = {};
  }

  render() {
    return (
      <div>
        <ConnectionConfig {...this.props} />
      </div>
    );
  }
}
