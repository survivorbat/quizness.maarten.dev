import Creator from "../models/creator";
import {Quiz} from "../models/quiz";
import {baseUrl} from "./constants";
import CreatorGameClient from "./creator-game-client";
import PlayerGameClient, {GameCallbacks} from "./player-game-client";
import Game from "../models/game";
import CreateGame from "../models/create-game";
import Player from "../models/player";

class BackendSdk {
  constructor(private readonly sdkToken: string | null) {
  }

  // Code is the OAuth authentication code to exchange with our own backend
  async authenticate(code: string): Promise<string> {
    const result = await fetch(`${baseUrl}/api/v1/tokens`, {method: 'post', body: JSON.stringify({code})});

    const token = result.headers.get('token')

    if (token) {
      return token
    }

    throw new Error('Failed to authenticate');
  }

  async refresh(): Promise<string> {
    const result = await fetch(`${baseUrl}/api/v1/tokens`, {method: 'put', headers: this.authHeader()});
    const token = result.headers.get('token')

    if (token) {
      return token
    }

    throw new Error('Failed to authenticate');
  }

  async getCreator(): Promise<Creator> {
    const result = await fetch(`${baseUrl}/api/v1/creators/self`, {headers: this.authHeader()});
    if (!result.ok) {
      throw new Error('Failed to fetch data');
    }

    return result.json();
  }

  async getQuizzes(): Promise<Quiz[]> {
    const result = await fetch(`${baseUrl}/api/v1/quizzes`, {headers: this.authHeader()});
    if (!result.ok) {
      throw new Error('Failed to fetch data');
    }

    return result.json();
  }

  async getGamesByQuiz(quiz: string): Promise<Game> {
    const result = await fetch(`${baseUrl}/api/v1/quizzes/${quiz}/games`, {headers: this.authHeader()});
    if (!result.ok) {
      throw new Error('Failed to fetch data');
    }

    const games: Game[] = await result.json();
    if (games.length === 1) {
      return games[1];
    }

    throw new Error('Game not found');
  }

  async getGameByCode(code: string): Promise<Game[]> {
    const result = await fetch(`${baseUrl}/api/v1/games?code=${code}`);
    if (!result.ok) {
      throw new Error('Failed to fetch data');
    }

    return result.json();
  }

  async createGame(quiz: string, game: CreateGame): Promise<void> {
    const data: RequestInit = {
      method: 'POST',
      headers: this.authHeader(),
      body: JSON.stringify(game),
    };

    const result = await fetch(`${baseUrl}/api/v1/quizzes/${quiz}/games`, data);
    if (!result.ok) {
      throw new Error('Failed to create game');
    }
  }

  async createPlayer(game: string): Promise<Player> {
    const data: RequestInit = {
      method: 'POST',
    };

    const result = await fetch(`${baseUrl}/api/v1/games/${game}/players`, data);
    if (!result.ok) {
      throw new Error('Failed to create game');
    }

    return result.json();
  }

  async getQuizByGame(game: string): Promise<Quiz> {
    const result = await fetch(`${baseUrl}/api/v1/games/${game}/quiz`);
    if (!result.ok) {
      throw new Error('Failed to fetch data');
    }

    return result.json();
  }

  getCreatorClient(game: string, callbacks: GameCallbacks): CreatorGameClient {
    return new CreatorGameClient(this.sdkToken!, game, callbacks);
  }

  getPlayerClient(game: string, player: string, callbacks: GameCallbacks): PlayerGameClient {
    return new PlayerGameClient(game, player, callbacks);
  }

  private authHeader(): Record<string, string> {
    return {'Authorization': `Bearer ${this.sdkToken}`}
  }
}

export default BackendSdk;