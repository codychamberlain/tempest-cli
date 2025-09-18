#!/usr/bin/env node

import { Command } from 'commander';
import dotenv from 'dotenv';
import axios from 'axios';


const program = new Command();
dotenv.config();

const api_token = process.env.API_TOKEN

function getData(url, station) {
  if (station === undefined) {
          axios.get(url, {
      params: {
          'token': api_token
      }
  })
      .then(response => {
          // Access the response data
          const responseData = response.data;
          console.log(responseData);
      })
      .catch(error => {
          console.error('Error fetching data:', error);
      });
  } else {
      axios.get(url, {
      params: {
          'token': api_token,
          'station_id': station
      }
  })
      .then(response => {
          // Access the response data
          const responseData = response.data;
          console.log(responseData);
      })
      .catch(error => {
          console.error('Error fetching data:', error);
      });
  }
  }


program
  .name('tempest-cli')
  .description('cli application utility for accessing data in tempestws API / weather stations')
  .version('1.0.0')
  .option('-s, --station <station>', 'id of the staiton to bind to for data');

program
  .command('forecast')
  .description('get forecast data for the station specified')
  .action((optons) => {
    const station_id = program.optsWithGlobals().station
    const forecast_url = 'https://swd.weatherflow.com/swd/rest/better_forecast'
    console.log(getData(forecast_url, station_id));
  });

  program
  .command('stations')
  .description('get station data for all stations or from a specific station if specified')
  .action((optons) => {
    const station_id = program.optsWithGlobals().station
    const forecast_url = 'https://swd.weatherflow.com/swd/rest/stations'
    console.log(getData(forecast_url, station_id));
  });

  program
  .command('observation')
  .description('get station data for all stations or from a specific station if specified')
  .action((optons) => {
    const station_id = program.optsWithGlobals().station
    const forecast_url = 'https://swd.weatherflow.com/swd/rest/observations/station/' + station_id
    console.log(getData(forecast_url, station_id));
  });

program.parse();