import { CustomVariableModel } from '@grafana/data';
import { getBackendSrv, setBackendSrv } from '@grafana/runtime';
import { TemplateSrv } from 'app/features/templating/template_srv';

import { ResourcesAPI } from '../resources/ResourcesAPI';

import { CloudWatchSettings, setupMockedTemplateService } from './CloudWatchDataSource';

export function setupMockedResourcesAPI({
  variables,
  response,
  getMock,
}: {
  getMock?: jest.Mock;
  response?: unknown;
  variables?: CustomVariableModel[];
  mockGetVariableName?: boolean;
} = {}) {
  let templateService = variables ? setupMockedTemplateService(variables) : new TemplateSrv();

  const api = new ResourcesAPI(CloudWatchSettings, templateService);
  let resourceRequestMock = getMock ? getMock : jest.fn().mockReturnValue(response);
  setBackendSrv({
    ...getBackendSrv(),
    get: resourceRequestMock,
  });

  return { api, resourceRequestMock, templateService };
}
