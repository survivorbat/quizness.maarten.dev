class BackendSdk {
    // May be empty
    private sdkToken?: string | null;

    constructor(private readonly baseUrl: string) {
        this.sdkToken = localStorage.getItem('token');
    }

    private authHeader(): Record<string, string> {
        return {'Authorization': `Bearer ${this.sdkToken}`}
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

        localStorage.removeItem('token');

        return false;
    }

    async refresh(): Promise<void> {
        const result = await fetch(`${this.baseUrl}/api/v1/tokens`, {method: 'put', headers: this.authHeader()});
        const token = result.headers.get('token')

        if (token) {
            localStorage.setItem('token', token)
            this.sdkToken = token;
            return
        }

        this.sdkToken = null;
        localStorage.removeItem('token');
    }
}

export default BackendSdk;