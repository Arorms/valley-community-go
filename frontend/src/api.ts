import axios from "axios";
import type {Post, PostCreateRequest} from "./types.ts";

const api = axios.create({
    baseURL: 'http://localhost:3000',
});

export const fetchPosts = async (offset = 0, limit = 10): Promise<Post[]> => {
    const response = await api.get<Post[]>(`/api/post/page?offset=${offset}&limit=${limit}`);
    return response.data;
};

export const fetchPostById = async (id: number): Promise<Post> => {
    const response = await api.get<Post>(`/api/post?id=${id}`);
    return response.data;
};

export const createPost = async (data: PostCreateRequest): Promise<Post> => {
    const response = await api.post<Post>('/api/post', data);
    return response.data;
};