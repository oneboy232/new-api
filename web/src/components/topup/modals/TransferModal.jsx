/*
Copyright (C) 2025 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/

import React, { useState, useEffect } from 'react';
import { Modal, Typography, Input, InputNumber } from '@douyinfe/semi-ui';
import { CreditCard } from 'lucide-react';
import { quotaToAmount, amountToQuota, getCurrencyConfig } from '../../../helpers/render';

const TransferModal = ({
  t,
  openTransfer,
  transfer,
  handleTransferCancel,
  userState,
  renderQuota,
  getQuotaPerUnit,
  transferAmount,
  setTransferAmount,
  title,
  availableLabel,
  amountLabel,
}) => {
  const isWithdraw = title?.includes('提现');
  const availableQuota = userState?.user?.aff_quota || 0;
  const { symbol } = getCurrencyConfig();

  // 金额输入模式（quota 或 钱）
  const [inputMode, setInputMode] = useState('money'); // 'quota' 或 'money'
  const [moneyAmount, setMoneyAmount] = useState(0);

  useEffect(() => {
    if (openTransfer) {
      setInputMode('money');
      setMoneyAmount(quotaToAmount(transferAmount));
    }
  }, [openTransfer]);

  const handleMoneyChange = (value) => {
    setMoneyAmount(value);
    setTransferAmount(amountToQuota(value));
  };

  const handleQuotaChange = (value) => {
    setTransferAmount(value);
    setMoneyAmount(quotaToAmount(value));
  };

  return (
    <Modal
      title={
        <div className='flex items-center'>
          <CreditCard className='mr-2' size={18} />
          {title || t('划转邀请额度')}
        </div>
      }
      visible={openTransfer}
      onOk={transfer}
      onCancel={handleTransferCancel}
      maskClosable={false}
      centered
    >
      <div className='space-y-4'>
        <div>
          <Typography.Text strong className='block mb-2'>
            {availableLabel || t('可用邀请额度')}
          </Typography.Text>
          <Input
            value={renderQuota(availableQuota)}
            disabled
            className='!rounded-lg'
          />
        </div>

        {/* 切换输入模式 */}
        <div className='flex items-center'>
          <div
            className='inline-flex items-center p-1 rounded-lg'
            style={{ backgroundColor: '#e6f4ff', border: '1px solid #91caff' }}
          >
            <button
              type='button'
              className={`px-4 py-1.5 text-sm font-medium rounded-md transition-all duration-200 ${
                inputMode === 'money'
                  ? 'shadow-sm'
                  : ''
              }`}
              style={
                inputMode === 'money'
                  ? { backgroundColor: '#1677ff', color: '#fff' }
                  : { color: '#1677ff' }
              }
              onClick={() => setInputMode('money')}
            >
              {t('金额')}
            </button>
            <button
              type='button'
              className={`px-4 py-1.5 text-sm font-medium rounded-md transition-all duration-200 ${
                inputMode === 'quota'
                  ? 'shadow-sm'
                  : ''
              }`}
              style={
                inputMode === 'quota'
                  ? { backgroundColor: '#1677ff', color: '#fff' }
                  : { color: '#1677ff' }
              }
              onClick={() => setInputMode('quota')}
            >
              {t('额度')}
            </button>
          </div>
        </div>

        <div>
          <Typography.Text strong className='block mb-2'>
            {amountLabel || t('划转额度')} · {t('最低') + renderQuota(getQuotaPerUnit())}
          </Typography.Text>
          {inputMode === 'quota' ? (
            <InputNumber
              min={getQuotaPerUnit()}
              max={availableQuota}
              value={transferAmount}
              onChange={handleQuotaChange}
              className='w-full !rounded-lg'
            />
          ) : (
            <InputNumber
              min={quotaToAmount(getQuotaPerUnit())}
              max={quotaToAmount(availableQuota)}
              value={moneyAmount}
              onChange={handleMoneyChange}
              formatter={(value) => `${symbol} ${value}`}
              parser={(value) => value.replace(new RegExp(`\\${symbol}\\s?`), '')}
              precision={2}
              className='w-full !rounded-lg'
            />
          )}
          <div className='text-xs text-gray-400 mt-1'>
            {inputMode === 'quota'
              ? `${t('约')} ${symbol}${quotaToAmount(transferAmount).toFixed(2)}`
              : `${t('约')} ${renderQuota(transferAmount)}`}
          </div>
        </div>
      </div>
    </Modal>
  );
};

export default TransferModal;
