import { GraphQLWsLink } from "@apollo/client/link/subscriptions";
import { createClient } from 'graphql-ws';
import { HttpLink, InMemoryCache } from '@apollo/client/core';
import { getMainDefinition } from "@apollo/client/utilities";
import { DocumentNode, Kind, OperationTypeNode } from "graphql";
import { ApolloLink } from "@apollo/client";
import { ApolloClient } from "@apollo/client";

// this is by far the cleanest aproach you can create than messy stuff
// inspired from angular apollo writing style
// env
const ragHTTP = 'http://localhost:8080/ask-chess-coach'
const ragWS = 'ws://localhost:8080/ask-chess-coach'
// end

// clients
const httpRagClient = new HttpLink({ uri: ragHTTP })
const gqlRagClient = new GraphQLWsLink(createClient({ url: ragWS }))
// end

// link
const ragLink = ApolloLink.split(({ query }: { query: DocumentNode }) => {

    const def = getMainDefinition(query)
    return def.kind === Kind.OPERATION_DEFINITION &&
        def.operation === OperationTypeNode.SUBSCRIPTION
},
    gqlRagClient,
    httpRagClient)
// end

export function InitGraphql() {
    return {
        rag: new ApolloClient({
            link: ragLink,
            cache: new InMemoryCache()
        })
    }
}

