import React, { useEffect, useState } from 'react';
import { fetchPosts, fetchPostById, createPost} from './api';
import type {Post} from "./types.ts";


export default function App() {
  const [posts, setPosts] = useState<Post[]>([]);
  const [selectedPost, setSelectedPost] = useState<Post | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [detailLoading, setDetailLoading] = useState<boolean>(false);

  // Form state
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');

  // Load the list on mount
  useEffect(() => {
    loadPostList();
  }, []);

  const loadPostList = async () => {
    setLoading(true);
    try {
      const data = await fetchPosts(0, 10);
      setPosts(data || []);
    } catch (err) {
      alert('Failed to load posts. Please check the backend service.');
    } finally {
      setLoading(false);
    }
  };

  // View details on click (calls GET /api/post?id=x)
  const handleSelectPost = async (id: number) => {
    setDetailLoading(true);
    try {
      const post = await fetchPostById(id);
      setSelectedPost(post);
    } catch (err) {
      alert('Failed to load post details');
    } finally {
      setDetailLoading(false);
    }
  };

  // Submit a new post (calls POST /api/post)
  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title || !content) return;

    try {
      await createPost({ title, content });
      setTitle('');
      setContent('');
      loadPostList(); // Refresh the list
    } catch (err) {
      alert('Failed to create post');
    }
  };

  return (
      <div className="min-h-screen bg-gray-100 p-8">
        <div className="max-w-5xl mx-auto grid grid-cols-1 md:grid-cols-2 gap-8">

          {/* Left: post form + list */}
          <div className="space-y-6">
            {/* Post form */}
            <div className="bg-white p-6 rounded-xl shadow-sm border border-gray-100">
              <h2 className="text-xl font-bold mb-4 text-gray-800">Create New Post</h2>
              <form onSubmit={handleCreate} className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Title</label>
                  <input
                      type="text"
                      value={title}
                      onChange={(e) => setTitle(e.target.value)}
                      className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                      placeholder="Enter a title"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Content</label>
                  <textarea
                      value={content}
                      onChange={(e) => setContent(e.target.value)}
                      className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 h-24"
                      placeholder="Enter content"
                  />
                </div>
                <button
                    type="submit"
                    className="w-full bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 rounded-lg transition"
                >
                  Publish
                </button>
              </form>
            </div>

            {/* Post list */}
            <div className="bg-white p-6 rounded-xl shadow-sm border border-gray-100">
              <h2 className="text-xl font-bold mb-4 text-gray-800">Posts</h2>
              {loading ? (
                  <p className="text-gray-500">Loading...</p>
              ) : posts.length === 0 ? (
                  <p className="text-gray-400">No posts yet</p>
              ) : (
                  <ul className="divide-y divide-gray-100">
                    {posts.map((post) => (
                        <li
                            key={post.id}
                            onClick={() => handleSelectPost(post.id)}
                            className={`py-3 px-2 cursor-pointer rounded-lg hover:bg-blue-50 transition ${
                                selectedPost?.id === post.id ? 'bg-blue-50 font-semibold' : ''
                            }`}
                        >
                          <div className="text-gray-900">{post.title}</div>
                          <div className="text-xs text-gray-400 mt-1">ID: {post.id}</div>
                        </li>
                    ))}
                  </ul>
              )}
            </div>
          </div>

          {/* Right: post details card */}
          <div className="bg-white p-6 rounded-xl shadow-sm border border-gray-100 h-fit sticky top-8">
            <h2 className="text-xl font-bold mb-4 text-gray-800 border-b pb-3">Post Details</h2>

            {detailLoading ? (
                <div className="text-gray-500 py-8 text-center">Loading details...</div>
            ) : selectedPost ? (
                <div className="space-y-4">
                  <div>
                <span className="inline-block bg-blue-100 text-blue-800 text-xs px-2.5 py-0.5 rounded font-medium mb-2">
                  ID: #{selectedPost.id}
                </span>
                    <h1 className="text-2xl font-bold text-gray-900">{selectedPost.title}</h1>
                  </div>
                  <hr className="border-gray-100" />
                  <div className="text-gray-700 whitespace-pre-line leading-relaxed">
                    {selectedPost.content}
                  </div>
                </div>
            ) : (
                <div className="text-gray-400 py-12 text-center">
                  👈 Click a post on the left to view details
                </div>
            )}
          </div>

        </div>
      </div>
  );
}