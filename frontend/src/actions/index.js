export const addTransaction = (trans) => {
  return {
    type: 'ADD_TRANSACTION',
    payload: trans
  }
}