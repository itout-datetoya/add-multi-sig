import React, { useState } from 'react';
import axios from 'axios';
import { connectMetamask, signWithMetamask } from '../services/metamask';

const API_URL = 'http://localhost:8080/api';

const Login: React.FC = () => {
  const [address, setAddress] = useState<string>('');
  const [loggedIn, setLoggedIn] = useState<boolean>(false);
  const [message, setMessage] = useState<string>('');

  const handleLogin = async () => {
    try {
      // 1. Metamaskからアカウント取得
      const addr = await connectMetamask();
      setAddress(addr);

      // 2. サーバーからチャレンジ（nonce）を取得
      const challengeRes = await axios.get(`${API_URL}/auth/challenge`, { params: { address: addr } });
      const challenge: string = challengeRes.data.challenge;
      
      // 3. Metamaskでチャレンジを署名
      const signature = await signWithMetamask(challenge);

      // 4. サーバーへ署名付きでログインリクエスト
      const loginRes = await axios.post(`${API_URL}/auth/login`, {
        address: addr,
        signature: signature,
      });

      if (loginRes.data.success) {
        setLoggedIn(true);
        setMessage('Logged in successfully');
      } else {
        setMessage('Login failed: ' + loginRes.data.message);
      }
    } catch (error: any) {
      console.error(error);
      setMessage('Login error: ' + error.message);
    }
  };

  return (
    <div>
      <h2>Login</h2>
      {loggedIn ? (
        <p>Logged in as: {address}</p>
      ) : (
        <button onClick={handleLogin}>Connect with Metamask</button>
      )}
      {message && <p>{message}</p>}
    </div>
  );
};

export default Login;
