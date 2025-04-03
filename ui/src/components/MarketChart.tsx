import React, { useState, useEffect } from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { MarketData } from '../types';

interface MarketChartProps {
  data?: MarketData[];
  title: string;
}

// Function to generate chart data based on MarketData
const generateChartData = (assets: MarketData[]) => {
  if (!assets || assets.length === 0) return [];

  // Create 10 time points for each asset
  return assets.map(asset => {
    const basePrice = asset.price || 100;
    const min = asset.low || basePrice * 0.95;
    const max = asset.high || basePrice * 1.05;
    
    // Calculate range for random fluctuation
    const range = max - min;
    
    return {
      name: asset.symbol,
      data: Array.from({ length: 10 }, (_, i) => {
        // More recent data points have higher weight towards current price
        const weight = i / 10;
        const randomValue = min + Math.random() * range;
        const weightedPrice = (weight * basePrice) + ((1 - weight) * randomValue);
        
        return {
          time: i,
          price: weightedPrice
        };
      })
    };
  });
};

const MarketChart: React.FC<MarketChartProps> = ({ data, title }) => {
  const [chartData, setChartData] = useState<any[]>([]);
  const [loading, setLoading] = useState<boolean>(true);

  // Initialize chart data
  useEffect(() => {
    if (!data || data.length === 0) {
      setChartData([]);
      setLoading(false);
      return;
    }

    setChartData(generateChartData(data));
    setLoading(false);
  }, [data]);
  
  // Simulate real-time updates
  useEffect(() => {
    if (!data || data.length === 0) return;
    
    const interval = setInterval(() => {
      setChartData(prevData => {
        return prevData.map(asset => {
          // Update each asset's data with a new point
          const newData = [...asset.data.slice(1)];
          const lastPoint = newData[newData.length - 1];
          
          // Generate a new point with small random movement
          const basePrice = lastPoint ? lastPoint.price : (data.find(d => d.symbol === asset.name)?.price || 100);
          const randomChange = (Math.random() - 0.5) * 0.02 * basePrice;
          
          newData.push({
            time: lastPoint ? lastPoint.time + 1 : 10,
            price: basePrice + randomChange
          });
          
          return {
            ...asset,
            data: newData
          };
        });
      });
    }, 5000);
    
    return () => clearInterval(interval);
  }, [data]);

  if (loading) {
    return (
      <div className="bg-white p-4 rounded-lg shadow-md">
        <h3 className="text-lg font-semibold mb-3">{title}</h3>
        <div className="h-60 flex items-center justify-center">
          <div className="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-blue-500"></div>
        </div>
      </div>
    );
  }

  // Pick colors for different lines
  const colors = ["#3b82f6", "#ef4444", "#22c55e", "#f97316", "#8b5cf6", "#06b6d4"];

  return (
    <div className="bg-white p-4 rounded-lg shadow-md">
      <h3 className="text-lg font-semibold mb-3">{title}</h3>
      <div className="h-60">
        {chartData.length > 0 ? (
          <ResponsiveContainer width="100%" height="100%">
            <LineChart margin={{ top: 5, right: 5, left: 5, bottom: 5 }}>
              <CartesianGrid strokeDasharray="3 3" opacity={0.2} />
              <XAxis dataKey="time" hide />
              <YAxis domain={['auto', 'auto']} hide />
              <Tooltip 
                formatter={(value: number) => [`$${value.toFixed(2)}`, 'Price']}
                labelFormatter={() => ''} // Hide the time label
              />
              {chartData.map((asset, index) => (
                <Line 
                  key={asset.name}
                  data={asset.data}
                  type="monotone"
                  dataKey="price"
                  name={asset.name}
                  stroke={colors[index % colors.length]}
                  dot={false}
                  activeDot={{ r: 5 }}
                />
              ))}
            </LineChart>
          </ResponsiveContainer>
        ) : (
          <div className="h-full flex items-center justify-center bg-gray-50 rounded border border-gray-200">
            <p className="text-center text-gray-500">No data available</p>
          </div>
        )}
      </div>
      
      {/* Legend */}
      {chartData.length > 0 && (
        <div className="mt-2 flex flex-wrap gap-3">
          {chartData.map((asset, index) => (
            <div key={asset.name} className="flex items-center text-sm">
              <div 
                className="w-3 h-3 mr-1 rounded-full" 
                style={{ backgroundColor: colors[index % colors.length] }}
              ></div>
              <span>{asset.name}</span>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default MarketChart; 