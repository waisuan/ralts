import React, { useEffect, useState } from "react"
import axios from 'axios';

const NewsFeed: React.FC = () => {
    const [data, setData] = useState([]);

    useEffect(() => {
        axios.get('http://localhost:8001/api/news-feed', { auth: { username: "username", password: "password" } })
            .then(response => {
                setData(response.data);
            })
            .catch(error => {
                console.error(error);
            });
      }, []);

    console.log(data);

    return (
        <div>
            {data.map((d) => (
                <p key={d['id']}>{d['title']}</p>
            ))}
        </div>
    );
}

export default NewsFeed