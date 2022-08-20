import React from 'react';

const App: React.FC = () => {

  const onClickToRegisterUser = React.useCallback(() => {
    console.log('onClickToRegisterUser')
  },[])

  const onClickToSendMessageToMe = React.useCallback(() => {
    console.log('onClickToSendMessageToMe')
  },[])

  return (
    <div className="App">
      <h1>푸쉬 사용자로 등록되기</h1>
      <div>
        <button onClick={() => { onClickToRegisterUser() }}>
          푸쉬 사용자로 등록되기
        </button>
      </div>
      <h1>나에게 푸쉬 전송하기</h1>
      <div>
        <button onClick={() => { onClickToSendMessageToMe() }}>
          나에게 푸쉬 전송하기
        </button>
      </div>
    </div>
  );
}

export default App;
