import { gql } from '@apollo/client'

export const GET_NOTIFICATIONS = gql`
    query GetNotes($userID: String!) {
        getNotes(userID: $userID) {
            id
            userID
            details
            timestamp
        }
    }
`

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
