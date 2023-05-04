import React, { useEffect } from 'react'

import { Observable, of } from 'rxjs'
import { map } from 'rxjs/operators'

import { createAggregateError } from '@sourcegraph/common'
import { gql } from '@sourcegraph/http-client'
import { H2, Text } from '@sourcegraph/wildcard'

import { queryGraphQL } from '../../../../backend/graphql'
import { FilteredConnection } from '../../../../components/FilteredConnection'
import { PageTitle } from '../../../../components/PageTitle'
import {
    DotComProductLicensesResult,
    DotComProductLicensesVariables,
    ProductLicenseFields,
} from '../../../../graphql-operations'
import { eventLogger } from '../../../../tracking/eventLogger'

import {
    siteAdminProductLicenseFragment,
    SiteAdminProductLicenseNode,
    SiteAdminProductLicenseNodeProps,
} from './SiteAdminProductLicenseNode'

interface Props {}

/**
 * Displays the product licenses that have been created on Sourcegraph.com.
 */
export const SiteAdminProductLicensesPage: React.FunctionComponent<React.PropsWithChildren<Props>> = () => {
    useEffect(() => eventLogger.logViewEvent('SiteAdminProductLicenses'), [])

    const nodeProps: Pick<SiteAdminProductLicenseNodeProps, 'showSubscription'> = {
        showSubscription: true,
    }

    return (
        <div className="site-admin-product-subscriptions-page">
            <PageTitle title="Product subscriptions" />
            <H2>License key lookup</H2>
            <Text>Find matching licenses and their associated product subscriptions.</Text>
            <FilteredConnection<ProductLicenseFields, Pick<SiteAdminProductLicenseNodeProps, 'showSubscription'>>
                className="list-group list-group-flush mt-3"
                noun="product license"
                pluralNoun="product licenses"
                queryConnection={queryLicenses}
                nodeComponent={SiteAdminProductLicenseNode}
                nodeComponentProps={nodeProps}
                emptyElement={<span className="text-muted mt-2">Enter a partial license key to find matches.</span>}
                noSummaryIfAllNodesVisible={true}
                autoFocus={true}
            />
        </div>
    )
}

function queryLicenses(args: {
    first?: number
    query?: string
}): Observable<DotComProductLicensesResult['dotcom']['productLicenses']> {
    const variables: Partial<DotComProductLicensesVariables> = {
        first: args.first,
        licenseKeySubstring: args.query,
    }
    return args.query
        ? queryGraphQL<DotComProductLicensesResult>(
              gql`
                  query DotComProductLicenses($first: Int, $licenseKeySubstring: String) {
                      dotcom {
                          productLicenses(first: $first, licenseKeySubstring: $licenseKeySubstring) {
                              nodes {
                                  ...ProductLicenseFields
                              }
                              totalCount
                              pageInfo {
                                  hasNextPage
                              }
                          }
                      }
                  }
                  ${siteAdminProductLicenseFragment}
              `,
              variables
          ).pipe(
              map(({ data, errors }) => {
                  if (!data?.dotcom?.productLicenses || (errors && errors.length > 0)) {
                      throw createAggregateError(errors)
                  }
                  return data.dotcom.productLicenses
              })
          )
        : of({
              __typename: 'ProductLicenseConnection' as const,
              nodes: [],
              totalCount: 0,
              pageInfo: { __typename: 'PageInfo' as const, hasNextPage: false, endCursor: null },
          })
}