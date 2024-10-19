"use client"

import React, { useState, ChangeEvent, FormEvent } from 'react';
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { AlertCircle, CheckCircle } from 'lucide-react';

// Define the Driver type
interface Driver {
  age: number | '';
  yearsOfExperience: number | '';
  hasValidLicense: boolean; // Make sure this is boolean
  vehicleAge: number | '';
  insuranceScore: number | '';
}

// Define the API response type
interface ApiResponse {
  results: {
    isEligible: boolean;
  };
}

// Function to call the backend API for eligibility evaluation
const evaluateEligibility = async (driver: { rule_ids: string[]; user_data: Driver }): Promise<boolean | null> => {
  try {
    const response = await fetch('http://localhost:8080/api/evaluate_rule', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(driver),
    });

    if (!response.ok) {
      throw new Error('Network response was not ok');
    }

    const data: ApiResponse = await response.json();
    return data.results.isEligible;
  } catch (error) {
    console.error('Error:', error);
    return null;
  }
};

export default function Home() {
  const [driver, setDriver] = useState<Driver>({
    age: '',
    yearsOfExperience: '',
    hasValidLicense: false,
    vehicleAge: '',
    insuranceScore: '',
  });
  const [result, setResult] = useState<boolean | null>(null);

  const handleInputChange = (e: ChangeEvent<HTMLInputElement>) => {
    const { name, value, type, checked } = e.target;
    setDriver(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value
    }));
  };

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const isEligible = await evaluateEligibility({
      rule_ids: ['30'], // Replace with the relevant rule IDs
      user_data: {
        age: parseInt(driver.age as string, 10),
        yearsOfExperience: parseInt(driver.yearsOfExperience as string, 10),
        vehicleAge: parseInt(driver.vehicleAge as string, 10),
        insuranceScore: parseInt(driver.insuranceScore as string, 10),
        hasValidLicense: driver.hasValidLicense,
      }
    });
    setResult(isEligible);
  };

  return (
    <div className="min-h-screen bg-gradient-to-b from-blue-50 to-white p-6">
      <Card className="max-w-2xl mx-auto">
        <CardHeader>
          <CardTitle className="text-2xl font-bold text-center">Ride-Sharing Driver Eligibility Form</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <Label htmlFor="age">Age</Label>
                <Input
                  id="age"
                  name="age"
                  type="number"
                  required
                  value={driver.age}
                  onChange={handleInputChange}
                  className="mt-1"
                />
              </div>
              <div>
                <Label htmlFor="yearsOfExperience">Years of Driving Experience</Label>
                <Input
                  id="yearsOfExperience"
                  name="yearsOfExperience"
                  type="number"
                  required
                  value={driver.yearsOfExperience}
                  onChange={handleInputChange}
                  className="mt-1"
                />
              </div>
              <div>
                <Label htmlFor="vehicleAge">Vehicle Age (years)</Label>
                <Input
                  id="vehicleAge"
                  name="vehicleAge"
                  type="number"
                  required
                  value={driver.vehicleAge}
                  onChange={handleInputChange}
                  className="mt-1"
                />
              </div>
              <div>
                <Label htmlFor="insuranceScore">Insurance Score</Label>
                <Input
                  id="insuranceScore"
                  name="insuranceScore"
                  type="number"
                  required
                  value={driver.insuranceScore}
                  onChange={handleInputChange}
                  className="mt-1"
                />
              </div>
              <div className="flex items-center space-x-2 mt-4">
                <Checkbox
                  id="hasValidLicense"
                  name="hasValidLicense"
                  checked={driver.hasValidLicense}
                  onCheckedChange={(checked: boolean) => setDriver(prev => ({ ...prev, hasValidLicense: checked }))} 
                />
                <Label htmlFor="hasValidLicense">I have a valid driver's license</Label>
              </div>

            </div>
            <Button type="submit" className="w-full">Check Eligibility</Button>
          </form>

          {result !== null && (
            <div className={`mt-6 p-4 rounded-md ${result ? 'bg-green-100' : 'bg-red-100'}`}>
              {result ? (
                <div className="flex items-center text-green-700">
                  <CheckCircle className="mr-2" />
                  <span>Congratulations! You are eligible to become a ride-sharing driver.</span>
                </div>
              ) : (
                <div className="flex items-center text-red-700">
                  <AlertCircle className="mr-2" />
                  <span>We're sorry, but you do not meet the eligibility criteria at this time.</span>
                </div>
              )}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
