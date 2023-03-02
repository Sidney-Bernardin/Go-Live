export const setSessionID = (sID) => localStorage.setItem('session_id', sID)
export const getSessionID = () => localStorage.getItem('session_id')
export const removeSessionID = () => localStorage.removeItem('session_id')
