export const baseUrl = process.env.REACT_APP_BACKEND_URL as string;

export const baseSocketUrl = (process.env.REACT_APP_BACKEND_URL as string).replace('http', 'ws');
