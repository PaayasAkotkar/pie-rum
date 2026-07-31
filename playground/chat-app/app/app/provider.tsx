'use client'
import { ApolloProvider } from "@apollo/client/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useRef, useState } from "react";
import { InitGraphql } from "./misc/graphql-config";
import { ResponsiveProvider } from "./services/cssx/responsive/services/responsive/use-responsive";

export default function Provider({ children }: Readonly<{ children: React.ReactNode }>) {
    //const [queryClient] = useState(() => new QueryClient({
    //    defaultOptions: {
    //        queries: {
    //            staleTime: 60 * 1000,
    //        },
    //    },
    //}));
    const client = InitGraphql().rag
    return (
        <>
            <ApolloProvider client={client}>
                    <ResponsiveProvider>
                        {children}
                    </ResponsiveProvider>
            </ApolloProvider>
        </>
    )
}