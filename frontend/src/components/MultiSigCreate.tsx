import React, { useState } from 'react';
import { createMultiSig } from '../services/api';

const MultiSigCreate: React.FC = () => {
  const [owner, setOwner] = useState<string>(''); // ログイン済みユーザーのEthereumアドレス
  const [participant1, setParticipant1] = useState<string>('');
  const [participant2, setParticipant2] = useState<string>('');
  const [message, setMessage] = useState<string>('');

  const handleCreate = async () => {
    // 例として、ownerはログイン済みのaddressとする
    const result = await createMultiSig({ owner, participants: [participant1, participant2] });
    setMessage(result.message);
  };

  return (
    <div>
      <h2>Create MultiSig</h2>
      <input
        type="text"
        placeholder="Your Ethereum Address (Owner)"
        value={owner}
        onChange={(e) => setOwner(e.target.value)}
      />
      <input
        type="text"
        placeholder="Participant 1 Ethereum Address"
        value={participant1}
        onChange={(e) => setParticipant1(e.target.value)}
      />
      <input
        type="text"
        placeholder="Participant 2 Ethereum Address"
        value={participant2}
        onChange={(e) => setParticipant2(e.target.value)}
      />
      <button onClick={handleCreate}>Create</button>
      {message && <p>{message}</p>}
    </div>
  );
};

export default MultiSigCreate;
