import { View, Text } from 'react-native'
import { SafeAreaView } from 'react-native-safe-area-context'
import { colors } from '@/src/styles/colors'

export default function ReposScreen() {
    return (
        <SafeAreaView style={{ flex: 1, backgroundColor: colors.background }}>
            <View style={{ flex: 1, justifyContent: 'center', alignItems: 'center' }}>
                <Text style={{ color: colors.textPrimary }}>Repos — coming soon</Text>
            </View>
        </SafeAreaView>
    )
}