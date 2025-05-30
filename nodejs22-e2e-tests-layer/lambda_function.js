// Import all dependencies
const AWS = require('aws-sdk');
const axios = require('axios');
const _ = require('lodash');
const playwright = require('playwright-core');
const puppeteer = require('puppeteer-core');
const express = require('express');
const chromium = require('@sparticuz/chromium');
const nodemon = require('nodemon');
const StreamZip = require('node-stream-zip');
const FormData = require('form-data');
const cors = require('cors');
const archiver = require('archiver');
const { EventBridgeClient } = require('@aws-sdk/client-eventbridge');

// Simple test for each dependency
async function testDependencies() {
    try {
        // Test AWS SDK
        const eventBridge = new EventBridgeClient({ region: 'us-east-1' });
        console.log('AWS SDK and EventBridge client initialized');

        // Test axios
        const axiosResponse = await axios.get('https://httpbin.org/get');
        console.log('Axios test successful:', axiosResponse.status);

        // Test lodash
        console.log('Lodash test:', _.chunk([1, 2, 3, 4], 2));

        // Test browser automation (Chromium)
        console.log('Chromium path:', chromium.path);
        const browser = await puppeteer.launch({
            args: chromium.args,
            executablePath: await chromium.executablePath(),
            headless: true
        });
        console.log('Puppeteer browser launched');
        await browser.close();

        // Test other dependencies
        const form = new FormData();
        form.append('field', 'value');
        console.log('FormData created');

        const zip = new StreamZip.async({ file: 'test.zip' });
        console.log('StreamZip initialized');

        console.log('All dependencies loaded successfully');
        return true;
    } catch (error) {
        console.error('Dependency test failed:', error);
        return false;
    }
}

exports.handler = async (event) => {
    console.log('Received event:', JSON.stringify(event, null, 2));

    // Test all dependencies
    const dependenciesOK = await testDependencies();
    const developer = process.env.DEVELOPER || 'unknown';

    return {
        statusCode: 200,
        body: JSON.stringify({
            message: `Hello from Lambda! Developer: ${developer}`,
            dependenciesTest: dependenciesOK ? 'PASSED' : 'FAILED',
            input: event,
            env: process.env.DEVELOPER,
            chromiumPath: chromium.path
        }),
    };
};