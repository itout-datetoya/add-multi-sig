import axios from 'axios';

const API_URL = 'http://localhost:8080/api';

interface CreateMultiSigData {
  owner: string;
  participants: string[];
}

export async function createMultiSig(data: CreateMultiSigData) {
  try {
    const res = await axios.post(`${API_URL}/multisig/create`, data);
    return res.data;
  } catch (error) {
    console.error(error);
    return { message: 'Error creating MultiSig' };
  }
}

export async function getMultiSigList() {
  try {
    const res = await axios.get(`${API_URL}/multisig/list`);
    return res.data;
  } catch (error) {
    console.error(error);
    return [];
  }
}

export async function signMultiSigData(id: string, method: 'get' | 'post', payload?: any) {
  try {
    if (method === 'get') {
      const res = await axios.get(`${API_URL}/multisig/${id}/data`);
      return res.data.dataToSign;
    } else {
      const res = await axios.post(`${API_URL}/multisig/${id}/data`, payload);
      return res.data;
    }
  } catch (error) {
    console.error(error);
    return { message: 'Error during signing process' };
  }
}
