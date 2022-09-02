/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { CancelablePromise } from '../core/CancelablePromise';
import type { BaseHttpRequest } from '../core/BaseHttpRequest';

export class InfoService {

  constructor(public readonly httpRequest: BaseHttpRequest) {}

  /**
   * GeoIP database info
   * Query the actively used GeoIP database info for all enabled database types.
   * @returns any OK
   * @throws ApiError
   */
  public getDatabaseMetadata({
    pretty = false,
  }: {
    pretty?: boolean,
  }): CancelablePromise<any> {
    return this.httpRequest.request({
      method: 'GET',
      url: '/metadata',
      query: {
        'pretty': pretty,
      },
    });
  }

  /**
   * Get OpenAPI spec
   * Returns the currently setup OpenAPI spec for this version of GeoIP. Note that some legacy/deprecated API endpoints are not included in this spec.
   * @returns any OK
   * @throws ApiError
   */
  public getOpenApiSpec(): CancelablePromise<any> {
    return this.httpRequest.request({
      method: 'GET',
      url: '/openapi.yaml',
    });
  }

  /**
   * Health check
   * Health check. Can also be used to check the rate-limit status without incrementing the rate-limit counters.
   * @returns any OK
   * @throws ApiError
   */
  public checkHealth(): CancelablePromise<{
    pong: boolean;
  }> {
    return this.httpRequest.request({
      method: 'GET',
      url: '/ping',
    });
  }

}
