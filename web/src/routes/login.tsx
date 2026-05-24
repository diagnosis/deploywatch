// src/routes/login.tsx
import { createFileRoute } from '@tanstack/react-router'
import {FaGithub, FaBell, FaBolt, FaCodeBranch, FaLifeRing, FaStream, FaSign, FaSignal} from 'react-icons/fa'

export const Route = createFileRoute('/login')({
    component: LoginPage,
})

function LoginPage() {
    return (
        <div className="min-h-screen bg-[#0a0a0f] flex items-center justify-center relative overflow-hidden">

            {/* Background grid */}
            <div
                className="absolute inset-0 opacity-[0.03]"
                style={{
                    backgroundImage: `linear-gradient(#fff 1px, transparent 1px), linear-gradient(90deg, #fff 1px, transparent 1px)`,
                    backgroundSize: '48px 48px',
                }}
            />

            {/* Glow */}
            <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] rounded-full bg-indigo-600/10 blur-[120px] pointer-events-none" />

            {/* Card */}
            <div className="relative z-10 w-full max-w-sm mx-4">

                {/* Logo */}
                <div className="flex items-center gap-2 mb-10 justify-center">
                    <div className="w-8 h-8 rounded-lg bg-indigo-500 flex items-center justify-center">
                        <FaBolt size={16} className="text-white" fill="white" />
                    </div>
                    <span className="text-white font-semibold text-lg tracking-tight">deploywatch</span>
                </div>

                {/* Main card */}
                <div className="bg-white/[0.03] border border-white/[0.08] rounded-2xl p-8 backdrop-blur-sm">

                    <h1 className="text-white text-2xl font-semibold tracking-tight mb-1">
                        Welcome back
                    </h1>
                    <p className="text-white/40 text-sm mb-8">
                        Sign in to monitor your repos in real time
                    </p>

                    {/* Feature pills */}
                    <div className="flex flex-wrap gap-2 mb-8">
                        {[
                            { icon: FaCodeBranch, label: 'Branch events' },
                            { icon: FaBell, label: 'Live notifications' },
                            { icon: FaSignal, label: 'Real-time SSE' },
                        ].map(({ icon: Icon, label }) => (
                            <div
                                key={label}
                                className="flex items-center gap-1.5 px-3 py-1 rounded-full bg-white/[0.04] border border-white/[0.06] text-white/40 text-xs"
                            >
                                <Icon size={11} />
                                {label}
                            </div>
                        ))}
                    </div>

                    {/* CTA */}
                    <a
                        href="http://localhost:8080/api/auth/github/login"
                        className="flex items-center justify-center gap-2.5 w-full py-3 px-4 rounded-xl bg-white text-[#0a0a0f] font-medium text-sm hover:bg-white/90 transition-colors duration-150"
                    >
                        <FaGithub size={16} />
                        Continue with GitHub
                    </a>

                    <p className="text-white/20 text-xs text-center mt-4">
                        By signing in you agree to our terms of service
                    </p>
                </div>
            </div>
        </div>
    )
}