import { Transaction } from '../Transaction';

export const addTransaction = (trans: Transaction) => {
  return {
    type: 'ADD_TRANSACTION',
    payload: trans
  };
};