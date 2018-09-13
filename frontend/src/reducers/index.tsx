import { Transaction } from '../Transaction';

interface TAction {
  type: string;
  payload: Transaction;
}

const initialState = {
  transactions: [] as Transaction[]
};

export default (state = initialState, action: TAction) => {
  switch (action.type) {
    case 'ADD_TRANSACTION':
      // console.log("Add transaction! :", action.payload)
      state.transactions.push(action.payload);
      return Object.assign({}, state, {
        transactions: state.transactions
      });
    default:
      return state;
  }
};