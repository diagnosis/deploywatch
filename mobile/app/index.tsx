import {useEffect, useState} from "react";
import {tokenStore} from "@/src/services/api";
import {Redirect} from "expo-router";
import {ActivityIndicator, View} from "react-native";
import {colors} from "@/src/styles/colors";


export default function Index(){
    const [loading, setLoading] = useState(false);
    const [hasToken, setHashToken] = useState(false);

    useEffect(() => {
        checkAuth()
    }, []);

    async function checkAuth(){
        try{
            const token = await tokenStore.getAccessToken()
            setHashToken(!!token)
        }
        catch {
            setHashToken(false)
        }finally {
            setLoading(false)
        }
    }
    if (loading) {
        return (
            <View style={{ flex: 1, justifyContent: 'center', alignItems: 'center', backgroundColor: colors.background }}>
                <ActivityIndicator size="large" color={colors.primary} />
            </View>
        )
    }

    return hasToken
        ? <Redirect href="/(tabs)/feed" />
        : <Redirect href="/login" />
}