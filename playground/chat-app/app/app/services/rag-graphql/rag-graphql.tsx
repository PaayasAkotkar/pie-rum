"use client";
import { useMutation, useSubscription } from "@apollo/client/react";
import { ChessStudent_REQUEST, CONFIG, ON_ChessCoach_REPLY, } from "./operation";
import { useCallback, useMemo } from "react";
import { ChessCoachPayload, ChessStudentRequest, IConfigRequest, OnChessCoachReply } from "./types";
import { InitGraphql } from "@/app/misc/graphql-config";

export function useRag() {
  const [askChessCoachMutation, { loading: mutationLoading, error: mutationError }] = useMutation
    //<
    //{ askChessCoach: ChessCoachPayload},
    //{ input: ChessStudentRequest }
    //>
    (ChessStudent_REQUEST)

  const askAi = async (input: ChessStudentRequest) => {
    try {
      const result = await askChessCoachMutation({ variables: { input } });
      console.log("AI Request sent:", result.data?.askChessCoach);
      return result.data?.askChessCoach;
    } catch (err) {
      console.error("Error sending AI request:", err);
      throw err;
    }
  };

  const [getConfig, { data,loading,error}] = useMutation(CONFIG)

  const config = useCallback(async (c: IConfigRequest) => {
    try {
      const res = await getConfig({
        variables: {
          input: c
        }
      })
      return res.data?.config
    } catch (err) {
      console.error(`updating config error ${err}`)
      throw error

    }
  }, [])

  
  const subscribeToReplies = (
    input: ChessStudentRequest,
  ) => {
    console.log("Starting subscription for room:", input.id);
    return null;
  }

  return {
   que: {
    
    },
    sub: {
      subToResp:subscribeToReplies
    },
    mut: {
      updateConfig: config,
      askAI:askAi
    },

    //askAi,
    //subscribeToReplies,
    //mutationLoading,
    //mutationError,
  };
}

export function useRagSubscription(input: ChessStudentRequest) {
  
  const { data, loading, error } = useSubscription(ON_ChessCoach_REPLY, {
    variables: { input },

    onData: ({ data }) => {
      console.log("Subscription data received:", data);
    },
    onError: (err) => {
      console.error("Subscription error:", err);
    },
    onComplete: () => {
      console.log("Subscription completed");
    },
  });

  return {
    aiReply: data?.chessCoachReply,
    subscriptionLoading: loading,
    subscriptionError: error,
  };
}