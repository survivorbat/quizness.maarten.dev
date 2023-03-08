import Creator from "../models/creator";

class BackendSdk {
  constructor(private readonly baseUrl: string, private readonly sdkToken: string | null) {
  }

  // Code is the OAuth authentication code to exchange with our own backend
  async authenticate(code: string): Promise<string> {
    const result = await fetch(`${this.baseUrl}/api/v1/tokens`, {method: 'post', body: JSON.stringify({code})});

    const token = result.headers.get('token')

    if (token) {
      return token
    }

    throw new Error('Failed to authenticate');
  }

  async refresh(): Promise<string> {
    const result = await fetch(`${this.baseUrl}/api/v1/tokens`, {method: 'put', headers: this.authHeader()});
    const token = result.headers.get('token')

    if (token) {
      return token
    }

    throw new Error('Failed to authenticate');
  }

  async getCreator(): Promise<Creator> {
    const result = await fetch(`${this.baseUrl}/api/v1/creators/self`, {headers: this.authHeader()});
    if (!result.ok) {
      throw new Error('Failed to fetch data');
    }

    return result.json();
  }

  private authHeader(): Record<string, string> {
    return {'Authorization': `Bearer ${this.sdkToken}`}
  }
}

export default BackendSdk;