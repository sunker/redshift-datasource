import defaults from 'lodash/defaults';

import React, { PureComponent } from 'react';
import { CodeEditor, CodeEditorSuggestionItem, CodeEditorSuggestionItemKind, InfoBox, LegacyForms } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from './datasource';
import { defaultQuery, RedShiftDataSourceJsonData, RedshiftQuery } from './types';
import { getTemplateSrv } from '@grafana/runtime';

type Props = QueryEditorProps<DataSource, RedshiftQuery, RedShiftDataSourceJsonData>;

export class QueryEditor extends PureComponent<Props> {
  onRawSqlChange = (rawSql: string) => {
    const { onChange, query, onRunQuery } = this.props;

    onChange({ ...query, rawSql });

    onRunQuery();
  };

  onChange = (value: RedshiftQuery) => {
    this.props.onChange(value);
    this.props.onRunQuery();
  };

  getSuggestions = (): CodeEditorSuggestionItem[] => {
    const sugs: CodeEditorSuggestionItem[] = [
      {
        label: '$__timeFilter',
        kind: CodeEditorSuggestionItemKind.Method,
        detail: '(Macro)',
      },
      {
        label: '$__timeGroup',
        kind: CodeEditorSuggestionItemKind.Method,
        detail: '(Macro)',
      },
    ];

    const templateSrv = getTemplateSrv();
    if (templateSrv) {
      templateSrv.getVariables().forEach(variable => {
        const label = `{${variable.name}}`;
        let val = templateSrv!.replace(label);
        if (val === label) {
          val = '';
        }
        sugs.push({
          label,
          kind: CodeEditorSuggestionItemKind.Text,
          detail: `(Template Variable) ${val}`,
        });
      });
    }
    return sugs;
  };

  render() {
    const query = defaults(this.props.query, defaultQuery);
    const { rawSql } = query;

    return (
      <>
        <InfoBox>To save and re-run the query, press ctrl/cmd+S.</InfoBox>
        <CodeEditor
          height={'250px'}
          language="sql"
          value={rawSql || ''}
          onBlur={this.onRawSqlChange}
          onSave={this.onRawSqlChange}
          showMiniMap={false}
          showLineNumbers={true}
          getSuggestions={this.getSuggestions}
        />
      </>
    );
  }
}
