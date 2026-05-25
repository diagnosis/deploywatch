import {QueryClient, QueryClientProvider} from "@tanstack/react-query";
import {useState} from "react";
import {SafeAreaProvider} from "react-native-safe-area-context";
import {Stack} from "expo-router";


export default function  RootLayout(){
    const [queryClient] = useState(() => new QueryClient({
        defaultOptions: {
            queries: {
                staleTime: 1000 * 60 * 5,
                retry: 1,
            }
        }
    }))

    return (
        <SafeAreaProvider>
            <QueryClientProvider client={queryClient}>
                <Stack screenOptions={{headerShown: false}}>
                    <Stack.Screen name="index" />
                    <Stack.Screen name="login" />
                    <Stack.Screen name="(tabs)" />
                </Stack>
            </QueryClientProvider>
        </SafeAreaProvider>
    )
}