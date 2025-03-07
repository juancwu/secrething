import type {
	User,
	LoginRequest,
	RegisterRequest,
	LockTOTPRequest,
	RemoveTOTPRequest,
	SetupTOTPResponse,
	LockTOTPResponse,
} from "@/lib/auth";

export type AuthProviderState = {
	user: User | null;
	isLoading: boolean;
	isAuthenticated: boolean;
	isPartialAuth: boolean;
	login: (credentials: LoginRequest) => Promise<void>;
	register: (data: RegisterRequest) => Promise<void>;
	logout: () => void;
	checkSession: () => Promise<void>;
	setupTOTP: () => Promise<SetupTOTPResponse>;
	lockTOTP: (data: LockTOTPRequest) => Promise<LockTOTPResponse>;
	removeTOTP: (data: RemoveTOTPRequest) => Promise<void>;
	resendVerificationEmail: () => Promise<void>;
};
