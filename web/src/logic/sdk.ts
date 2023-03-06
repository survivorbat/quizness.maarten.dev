class BackendSdk {
    // May be empty
    private sdkToken?: string | null;

    constructor(private readonly baseUrl: string) {
        this.sdkToken = localStorage.getItem('token');
    }

    // Code is the OAuth authentication code to exchange with our own backend
    async authenticate(code: string): Promise<boolean> {
        const result = await fetch(`${this.baseUrl}/api/v1/tokens`, {method: 'post', body: JSON.stringify({code})});

        const token = result.headers.get('token')

        if (token) {
            localStorage.setItem('token', token)
            this.sdkToken = token;
            return true
        }

        return false;
    }
}

export default BackendSdk;