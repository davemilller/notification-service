import { gql } from '@apollo/client'

export const NOTIFICATIONS_SUBSCRIPTION = gql`
    subscription Notifications($userID: String!) {
        notifications(userID: $userID) {
            id
            userID
            details
            timestamp
        }
    }
`
