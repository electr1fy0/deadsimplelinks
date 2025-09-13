import { useState, useEffect } from "react";

const BookmarkApp = () => {
  const [token, setToken] = useState(localStorage.getItem("token") || "");
  const [username, setUsername] = useState(
    localStorage.getItem("username") || "",
  );
  const [bookmarks, setBookmarks] = useState(null);
  const [newBookmark, setNewBookmark] = useState("");
  const [loginInput, setLoginInput] = useState("");
  const [loggingIn, setLoggingIn] = useState(false);

  useEffect(() => {
    if (!token) return;
    setBookmarks(null);
    fetch(`/api/bookmarks_get?token=${token}`)
      .then((res) => res.ok && res.json())
      .then((data) => data && setBookmarks(data))
      .catch(() => {
        setToken("");
        setUsername("");
        localStorage.removeItem("token");
        localStorage.removeItem("username");
        setBookmarks([]);
      });
  }, [token]);

  const handleLogin = async (e) => {
    e.preventDefault();
    if (!loginInput.trim()) return;

    setLoggingIn(true);

    const res = await fetch(`/api/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username: loginInput }),
    });

    if (res.ok) {
      const data = await res.json();
      localStorage.setItem("token", data.token);
      localStorage.setItem("username", loginInput);
      setToken(data.token);
      setUsername(loginInput);
      setLoginInput("");
    }

    setLoggingIn(false);
  };

  const handleLogout = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("username");
    setToken("");
    setUsername("");
    setBookmarks([]);
  };

  const addBookmark = (e) => {
    e.preventDefault();
    if (!newBookmark.trim() || !token) return;

    setBookmarks((prev) => [...(prev || []), newBookmark]);
    const toPost = newBookmark;
    setNewBookmark("");

    fetch(`/api/bookmarks_post`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ token, url: toPost }),
    }).catch(() => {
      setBookmarks((prev) => prev.filter((b) => b !== toPost));
    });
  };

  return (
    <main className="flex flex-col justify-center items-center min-h-screen w-screen px-4 bg-white">
      <div className="w-full max-w-[540px] mx-auto relative">
        {loggingIn && (
          <div className="absolute top-0 left-0 h-[2px] bg-neutral-800 animate-[loading_1s_linear_infinite]"></div>
        )}

        <div className="w-full">
          <p className="text-lg text-black">Dead Simple Bookmarks</p>
          <p className="text-neutral-500 text-sm">No password bookmarks.</p>
        </div>

        {!token ? (
          <div className="mt-8 mb-4 w-full">
            <form onSubmit={handleLogin} className="w-full">
              <input
                type="text"
                placeholder="Enter username..."
                value={loginInput}
                onChange={(e) => setLoginInput(e.target.value)}
                disabled={loggingIn}
                className="bg-neutral-100 text-black py-2 rounded-lg w-full text-sm px-4 focus:outline-none focus:ring-0 placeholder:text-neutral-600 disabled:opacity-50"
              />
            </form>
          </div>
        ) : (
          <div className="mt-8 mb-4 w-full">
            <form onSubmit={addBookmark} className="w-full">
              <input
                type="text"
                placeholder="Add a new link..."
                value={newBookmark}
                onChange={(e) => setNewBookmark(e.target.value)}
                className="bg-neutral-100 text-black py-2 rounded-lg w-full text-sm px-4 focus:outline-none focus:ring-0 placeholder:text-neutral-600"
              />
            </form>

            {bookmarks === null ? (
              <p className="text-neutral-400 text-sm my-4">
                Loading bookmarks...
              </p>
            ) : (
              <ul className="text-neutral-600 my-4 space-y-2 text-sm">
                {bookmarks.map((b, i) => (
                  <li key={i}>
                    <a
                      href={b}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="hover:underline"
                    >
                      {b}
                    </a>
                  </li>
                ))}
              </ul>
            )}

            <button
              onClick={handleLogout}
              className="w-full bg-neutral-100 text-black py-2 rounded-lg text-sm hover:bg-neutral-200"
            >
              Logout ({username})
            </button>
          </div>
        )}
      </div>

      <style>{`
        @keyframes loading {
          0% { width: 0%; }
          100% { width: 100%; }
        }
      `}</style>
    </main>
  );
};

export default BookmarkApp;
