export async function connectMetamask(): Promise<string> {
    if (window.ethereum) {
      try {
        const accounts = await window.ethereum.request({ method: 'eth_requestAccounts' });
        return accounts[0];
      } catch (error) {
        console.error('User rejected connection', error);
        throw error;
      }
    } else {
      alert('Please install Metamask!');
      throw new Error('Metamask not found');
    }
  }
  
  export async function signWithMetamask(data: string): Promise<string> {
    if (window.ethereum) {
      try {
        const account = await connectMetamask();
        const signature = await window.ethereum.request({
          method: 'personal_sign',
          params: [data, account],
        });
        return signature;
      } catch (error) {
        console.error('Signature error', error);
        throw error;
      }
    } else {
      throw new Error('Metamask not found');
    }
  }
  