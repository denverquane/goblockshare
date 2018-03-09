const initialState = {
  transactions: []
}

export default (state = initialState, action) => {
  switch (action.type) {
    case 'ADD_TRANSACTION':
      console.log("Add transaction! :", action.payload)
      state.transactions.push(action.payload);
      return Object.assign({}, state, {
        transactions: state.transactions
      })
    default:
      return state
  }
}