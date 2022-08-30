/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Address } from '../models/Address';
import type { BulkGeoResult } from '../models/BulkGeoResult';
import type { GeoResult } from '../models/GeoResult';
import type { LookupOptions } from '../models/LookupOptions';

import type { CancelablePromise } from '../core/CancelablePromise';
import type { BaseHttpRequest } from '../core/BaseHttpRequest';

export class LookupService {

  constructor(public readonly httpRequest: BaseHttpRequest) {}

  /**
   * Bulk lookup addresses
   * Lookup addresses (hostnames/domains/ipv4/ipv6/etc) in bulk.
   * **NOTE: A maximum of 25 addresses can be looked up at once, anything above will return a `400`**.
   * @returns BulkGeoResult Bulk lookup response, including successful and failed lookups.
   * @throws ApiError
   */
  public getManyAddresses({
    requestBody,
  }: {
    requestBody: {
      /**
       * Array of addresses to lookup.
       */
      addresses: Array<Address>;
      options?: LookupOptions;
    },
  }): CancelablePromise<BulkGeoResult> {
    return this.httpRequest.request({
      method: 'POST',
      url: '/bulk',
      body: requestBody,
      mediaType: 'application/json',
    });
  }

  /**
   * Lookup address
   * Lookup an address (hostname/domain/ipv4/ipv6/etc), returning the GeoIP information if available.
   * @returns GeoResult Response was successful.
   * @throws ApiError
   */
  public getAddress({
    address,
    pretty = false,
    disableHostLookup = false,
  }: {
    address: Address,
    pretty?: boolean,
    disableHostLookup?: boolean,
  }): CancelablePromise<GeoResult> {
    return this.httpRequest.request({
      method: 'GET',
      url: '/lookup/{address}',
      path: {
        'address': address,
      },
      query: {
        'pretty': pretty,
        'disable_host_lookup': disableHostLookup,
      },
    });
  }

}
