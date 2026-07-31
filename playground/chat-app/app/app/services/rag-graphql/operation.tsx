import { gql, TypedDocumentNode } from "@apollo/client";
import { ChessStudentRequest, IConfigRequest, IConfigResp, OnChessCoachReply } from "./types";

export const INIT_CHESS_SUGESSTION = gql`
  query  initCoachSuggestions($input: ChessStudentRequest) {
     initCoachSuggestions(input: $input) {
      information {
        year
        title
        desc
        outro
        miscItems {
          key
          value {
             title
             desc
             copy
             link
             canCopy
             isLink
          }
        }
      }
      suggestion {
        year
        title
        desc
        outro
        miscItems {
          key
          value {
             title
             desc
             copy
             link
             canCopy
             isLink
          }
        }
      }
      bestPractice {
        year
        title
        desc
        outro
        miscItems {
          key
          value {
             title
             desc
             copy
             link
             canCopy
             isLink
          }
        }
      }
      miscItems {
        key
        value {
           title
           desc
           canCopy
           copy
           link
           isLink
        }
      }
    }
  }
`

export const ON_ChessCoach_REPLY: TypedDocumentNode<{chessCoachReply:OnChessCoachReply},{input:ChessStudentRequest}> = gql`
  subscription chessCoachReply($input: ChessStudentRequest) {
    chessCoachReply(input: $input) {
     information {
        year
        title
        desc
        outro
        miscItems {
          key
          value {
             title
             desc
             copy
             link
             canCopy
             isLink
          }
        }
      }
      suggestion {
        year
        title
        desc
        outro
        miscItems {
          key
          value {
             title
             desc
             canCopy
             copy
             link
             isLink
          }
        }
      }
      bestPractice {
        year
        title
        desc
        outro
        miscItems {
          key
          value {
             title
             desc
             canCopy
             copy
             link
             isLink
          }
        }
      }
      miscItems {
        key
        value {
           title
           desc
           canCopy
           copy
           link
           isLink
        }
      }
    }
  }
`

export const ChessStudent_REQUEST: TypedDocumentNode<{ askChessCoach: OnChessCoachReply }, { input: ChessStudentRequest }> = gql`
  mutation askChessCoach($input: ChessStudentRequest) {
    askChessCoach(input: $input) {
      status
      message
    }
  }
`
export const CONFIG: TypedDocumentNode<{ config: IConfigResp }, { input: IConfigRequest }> = gql`
mutation config($input: IConfigRequest){
  config(input:$input){
succeed 
error
  }
}
`