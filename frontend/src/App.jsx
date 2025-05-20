import { useState,useEffect  } from "react";
import api from "./api";
import "./app.css"

export default function CountdownApp() {
  const [events, setEvents] = useState([]);
  const [title, setTitle] = useState("");
  const [date, setDate] = useState("");

  useEffect(() => {
  api.get("/")
    .then((res) => setEvents(res.data))
    .catch((err) => console.error("Failed to fetch events:", err));
}, []);

  const addCountdown = () => {    
    if (!title || !date) return;

    const newEvent = { title, date };

    api.post("/", newEvent)
      .then((res) => {
        setEvents([...events, res.data]); // assume backend returns the saved event
        setTitle("");
        setDate("");
      })
      .catch((err) => console.error("Failed to save event:", err));
  };

  const calculateDaysLeft = (targetDate) => {
    const now = new Date();
    const target = new Date(targetDate);
    const diffTime = target - now;
    return Math.ceil(diffTime / (1000 * 60 * 60 * 24));
  };

  const deleteCountdown = (id) => {
    api.delete(`/${id}`)
      .then(() => {
        setEvents(events.filter((event) => event.id !== id));
      })
      .catch((err) => console.error("Failed to delete event:", err));
  };
  
  return (
    <div className="min-h-screen min-w-screen  bg-gray-900 p-6 text-center">
      <h1 className="text-3xl font-bold mb-6">📆 Countdown Tracker</h1>

      <div className="flex justify-center gap-4 mb-6">
        <input
          type="text"
          placeholder="Event Title"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          className="border p-2 rounded"
        />
        <input
          type="date"
          value={date}
          onChange={(e) => setDate(e.target.value)}
          className="border p-2 rounded"
        />
        <button
          onClick={addCountdown}
          className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
        >
          Add
        </button>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {events.map((event, index) => {
          const daysLeft = calculateDaysLeft(event.date);
          return (
            <div key={index} className="bg-gray-700 shadow-md p-4 rounded">
              <h2 className="text-xl font-semibold">{event.title}</h2>
              <p className="text-gray-400">
                {daysLeft > 0
                  ? `${daysLeft} day(s) left`
                  : "Date has passed!"}
              </p>
              <p className="text-sm text-gray-500">{event.date}</p>
              <button
                onClick={() => deleteCountdown(event.id)}
                className="mt-2 text-white px-3 py-1 rounded "
              >
                Delete
              </button>
            </div>
          );
        })}
      </div>

      <h1 className="text-white p5 font-semibold m-10">
        "The only path is forward"
      </h1>
    </div>
  );
}
