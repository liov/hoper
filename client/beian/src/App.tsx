import { useState } from 'react'
import reactLogo from './assets/react.svg'
import beian from './assets/beian.png'
import './App.css'

function App() {
  const [count, setCount] = useState(0)

  return (
    <div className="App">
      <div>
      <iframe src="https://m.hoper.xyz" height="812" width="375" frameBorder="0"></iframe>
      </div>
      <div className="card">
        <button onClick={() => setCount((count) => count + 1)}>
          count is {count}
        </button>
      </div>
      <p className="read-the-docs">
        <a href="https://beian.miit.gov.cn/" target="_blank">晋ICP备18012261号-1</a>
      </p>
      <div  className="beian">
        <a
          target="_blank"
          href="http://www.beian.gov.cn/portal/registerSystemInfo?recordcode=44030602007392"
        ><img src={beian} />
          <p>
            粤公网安备 44030602007392号
          </p></a
        >
      </div>
    </div>
  )
}

export default App
