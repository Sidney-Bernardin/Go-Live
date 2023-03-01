export const setSessionID = (sID) => sessionStorage.setItem('session_id', sID)
export const getSessionID = () => sessionStorage.getItem('session_id')
export const removeSessionID = () => sessionStorage.removeItem('session_id')
