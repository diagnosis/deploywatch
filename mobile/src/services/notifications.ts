import * as Notifications from 'expo-notifications'
import {apiClient} from "@/src/services/api";
import {Platform} from "react-native";


export async function registerForPushNotifications(){
    const { status } = await Notifications.requestPermissionsAsync();
    if (status !== "granted") return

    const token = await Notifications.getDevicePushTokenAsync();
    if (!token.data) return

    try {
        await apiClient.post('/api/device-tokens',{
            token: token.data,
            platform: Platform.OS,
        })
        console.log("device token registered successfully")
    }catch (e){
        console.log("failed to register device token")
    }
}