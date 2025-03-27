import React, { useState } from 'react';
import { signMultiSigData } from '../services/api';
import { signWithMetamask } from '../services/metamask';

const MultiSigSign: React.FC = () => {
  const [multiSigId, setMultiSigId] = useState<string>('');
  const [message, setMessage] = useState<string>('');

  const handleSign = async () => {
    // GETで署名用のデータを取得（プロトコルの状態に応じた処理）
    const dataToSign = await signMultiSigData(multiSigId, 'get');
    // Metamaskの署名機能で署名
    const signature = await signWithMetamask(dataToSign);
    // 署名をPOSTで送信して状態を更新
    const result = await signMultiSigData(multiSigId, 'post', { signature });
    setMessage(result.message);
  };

  return (
    <div>
      <h2>MultiSig Sign</h2>
      <input
        type="text"
        placeholder="MultiSig ID"
        value={multiSigId}
        onChange={(e) => setMultiSigId(e.target.value)}
      />
      <button onClick={handleSign}>Sign and Update</button>
      {message && <p>{message}</p>}
    </div>
  );
};

export default MultiSigSign;
