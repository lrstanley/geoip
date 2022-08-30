/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { BulkError } from './BulkError';
import type { GeoResult } from './GeoResult';

export type BulkGeoResult = {
  errors: Array<BulkError>;
  results: Array<GeoResult>;
};

