/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

export type LookupOptions = {
  /**
   * Disable doing a reverse lookup of the resolved IP (much faster).
   */
  disable_host_lookup?: boolean;
  /**
   * BCP47 or standard 2-character language code.
   */
  lang?: string;
  /**
   * Pretty print the JSON response.
   */
  pretty?: boolean;
};

