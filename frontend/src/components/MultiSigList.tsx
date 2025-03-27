import React, { useEffect, useState } from 'react';
import { getMultiSigList } from '../services/api';

interface MultiSig {
  id: string;
  owner: string;
  participants: string[];
  status: string;
}

const MultiSigList: React.FC = () => {
  const [multiSigs, setMultiSigs] = useState<MultiSig[]>([]);

  useEffect(() => {
    async function fetchData() {
      const list = await getMultiSigList();
      setMultiSigs(list);
    }
    fetchData();
  }, []);

  return (
    <div>
      <h2>Your MultiSig List</h2>
      <ul>
        {multiSigs.map((ms) => (
          <li key={ms.id}>
            ID: {ms.id} | Owner: {ms.owner} | Participants: {ms.participants.join(', ')} | Status: {ms.status}
          </li>
        ))}
      </ul>
    </div>
  );
};

export default MultiSigList;
