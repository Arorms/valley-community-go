export interface Post {
    id: number;
    title: string;
    content: string;
    user_id: string;
}

export interface PostCreateRequest {
    title: string;
    content: string;
}