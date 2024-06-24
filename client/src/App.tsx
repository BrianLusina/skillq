import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Home from './pages/Home';
import NewProgrammer from './pages/NewProgrammer';
import Header from './components/Header';

function App() {
  return (
    <Router>
      <div>
        <Header />
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/newprogrammer" element={<NewProgrammer />} />
        </Routes>
      </div>
    </Router>
  )
}

export default App;